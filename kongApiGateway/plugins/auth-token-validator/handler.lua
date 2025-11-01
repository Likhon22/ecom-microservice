local AuthTokenValidator = {
  PRIORITY = 1001,
  VERSION = "1.0.0",
}

function AuthTokenValidator:access(conf)
  local cookie_header = kong.request.get_header("cookie")
  if not cookie_header then
    return return_error(401, "No authentication cookies")
  end

  local access_token = extract_cookie(cookie_header, "access-token")
  if not access_token then
    return return_error(401, "Access token not found")
  end

  -- Verify and decode JWT
  local payload, err = verify_jwt(access_token, conf.jwt_secret)
  if err then
    kong.log.warn("JWT verification failed: ", err)
    return return_error(401, err)
  end

  -- Check expiration
  local now = ngx.time()
  if payload.exp and payload.exp < now then
    return return_error(401, "Access token expired")
  end

  kong.log.info("Authenticated user: ", payload.Email, " (", payload.Role, ")")
end

-- Verify JWT with HMAC-SHA256 signature
function verify_jwt(token, secret)
  local parts = {}
  for part in token:gmatch("[^.]+") do
    table.insert(parts, part)
  end
  
  if #parts ~= 3 then
    return nil, "Invalid JWT format"
  end

  local header_payload = parts[1] .. "." .. parts[2]
  local signature = parts[3]

  -- Use Kong's built-in OpenSSL HMAC
  local hmac = require "resty.openssl.hmac"
  local h = hmac.new(secret, "sha256")
  h:update(header_payload)
  local digest = h:final()

  -- Convert to base64url
  local expected_sig = ngx.encode_base64(digest)
  expected_sig = expected_sig:gsub('+', '-'):gsub('/', '_'):gsub('=', '')

  -- Compare signatures
  if signature ~= expected_sig then
    return nil, "Invalid signature"
  end

  -- Decode payload
  local payload_b64 = parts[2]
  local padding = #payload_b64 % 4
  if padding > 0 then
    payload_b64 = payload_b64 .. string.rep('=', 4 - padding)
  end
  payload_b64 = payload_b64:gsub('-', '+'):gsub('_', '/')
  
  local payload_json = ngx.decode_base64(payload_b64)
  if not payload_json then
    return nil, "Failed to decode payload"
  end
  
  local cjson = require "cjson.safe"
  local payload, err = cjson.decode(payload_json)
  if not payload then
    return nil, "Failed to parse JSON: " .. (err or "unknown")
  end
  
  return payload, nil
end

function extract_cookie(cookie_str, name)
  if not cookie_str then return nil end
  local pattern = name .. "=([^;]+)"
  return cookie_str:match(pattern)
end

function return_error(status, message)
  return kong.response.exit(status, {
    success = false,
    message = message,
    status_code = status
  })
end

return AuthTokenValidator