local cjson = require "cjson.safe"

local UserEmailInjector = {
  PRIORITY = 900,  
  VERSION = "1.0.0",
}

function UserEmailInjector:access(conf)
  -- Get cookie
  local cookie_header = kong.request.get_header("cookie")
  if not cookie_header then
    return  -- No cookies, let validator handle it
  end

  -- Extract token
  local access_token = extract_cookie(cookie_header, "access-token")
  if not access_token then
    return  -- No token, let validator handle it
  end

  -- Decode payload (NO signature verification - validator already did it)
  local payload = decode_jwt_payload(access_token)
  if not payload or not payload.Email then
    kong.log.warn("Failed to extract email from JWT")
    return
  end

  -- Set email header
  kong.service.request.set_header("x-user-email", payload.Email)
  kong.log.info("Injected user email: ", payload.Email)
end

-- Just decode, don't verify signature
function decode_jwt_payload(token)
  local parts = {}
  for part in token:gmatch("[^.]+") do
    table.insert(parts, part)
  end

  if #parts ~= 3 then return nil end

  -- Decode middle part (payload)
  local payload_b64 = parts[2]
  local padding = #payload_b64 % 4
  if padding > 0 then
    payload_b64 = payload_b64 .. string.rep('=', 4 - padding)
  end
  payload_b64 = payload_b64:gsub('-', '+'):gsub('_', '/')

  local payload_json = ngx.decode_base64(payload_b64)
  if not payload_json then return nil end

  return cjson.decode(payload_json)
end

function extract_cookie(cookie_str, name)
  if not cookie_str then return nil end
  for cookie in cookie_str:gmatch("[^;]+") do
    local k, v = cookie:match("^%s*(.-)%s*=%s*(.-)%s*$")
    if k == name then return v end
  end
  return nil
end

return UserEmailInjector
