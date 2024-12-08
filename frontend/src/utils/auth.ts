const TOKEN_KEY = 'token'

function isLogin() {
  return !!localStorage.getItem(TOKEN_KEY)
}

function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export { clearToken, getToken, isLogin, setToken }
