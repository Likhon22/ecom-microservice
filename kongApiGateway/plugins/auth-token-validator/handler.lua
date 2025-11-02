local cjson = require "cjson.safe"
local ngx = ngx

local AuthTokenValidator = {
  PRIORITY = 1001,
  VERSION = "1.0.0",
}

-- Main access function
function AuthTokenValidator:access(conf)
  kong.log.warn("=== AUTH PLUGIN DEBUG ===")

  -- 1 Get JWT secret from environment
  local secret = conf.jwt_secret
  if not secret or secret == "" then
    kong.log.err("JWT secret not configured")
    return kong.response.exit(500, { 
      success = false, 
      message = "JWT secret missing" 
    })
  end

  -- 2 Get cookie header
  local cookie_header = kong.request.get_header("cookie")
  if not cookie_header then
    return return_error(401, "No authentication cookies")
  end

  -- 3 Extract access token from cookie
  local access_token = extract_cookie(cookie_header, "access-token")
  if not access_token then
    return return_error(401, "Access token not found")
  end

  -- 4 Verify JWT
  local payload, err = verify_jwt(access_token, secret)
  if err then
    kong.log.warn("JWT verification failed: ", err)
    return return_error(401, err)
  end

  -- 5 Check expiration
  local now = ngx.time()
  if payload.exp and payload.exp < now then
    return return_error(401, "Access token expired")
  end

  kong.log.warn("Authenticated user: ", payload.Email, " (", payload.Role, ")")
end

-- Helper: verify JWT using HMAC-SHA256
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

  local hmac = require "resty.openssl.hmac"
  local h = hmac.new(secret, "sha256")
  h:update(header_payload)
  local digest = h:final()

  -- Convert to base64url
  local expected_sig = ngx.encode_base64(digest):gsub('+','-'):gsub('/','_'):gsub('=','')

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

  local payload, err = cjson.decode(payload_json)
  if not payload then
    return nil, "Failed to parse JSON: " .. (err or "unknown")
  end

  return payload, nil
end

-- Helper: extract cookie
function extract_cookie(cookie_str, name)
  if not cookie_str then return nil end
  for cookie in cookie_str:gmatch("[^;]+") do
    local k, v = cookie:match("^%s*(.-)%s*=%s*(.-)%s*$")
    if k == name then
      return v
    end
  end
  return nil
end

-- Helper: return error response
function return_error(status, message)
  return kong.response.exit(status, {
    success = false,
    message = message,
    status_code = status
  })
end

return AuthTokenValidator
