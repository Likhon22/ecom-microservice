local AuthCookieClearer = {
  PRIORITY = 999,
  VERSION = "1.0.0"
}
function AuthCookieClearer:header_filter(conf)
  local status=kong.response.get_status();
  if status == 200 then
    local clear_cookies = {
      "access-token=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict",
      "refresh-token=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"
    }
    kong.response.set_header("Set-Cookie",clear_cookies)
    kong.log.info("Cleared auth cookies after successful logout")

  end
  
end
return AuthCookieClearer