
local  AuthMetaDataSetter ={
  PRIORITY = 1002,
  VERSION = "1.0.0",
}

function AuthMetaDataSetter:access(conf)
  local cookie_header = kong.request.get_header("Cookie")
  
  if not cookie_header then
    return  -- Just continue, don't block
  end
  
  local refresh_token = extract_cookie(cookie_header, "refresh-token")  -- ✅ hyphen
  
  if refresh_token then
    kong.service.request.set_header("refresh-token", refresh_token)  -- ✅ hyphen
    kong.log.info("Forwarded refresh-token to gRPC service")
  end
end

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

function return_error(status, message)
  return kong.response.exit(status, {
    success = false,
    message = message,
    status_code = status
  })
end

return AuthMetaDataSetter