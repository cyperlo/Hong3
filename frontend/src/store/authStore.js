import { reactive } from 'vue';

const authState = reactive({
  user: null,
  token: null,
  isAuthenticated: false,
});

// 从localStorage加载认证信息
const loadAuth = () => {
  const token = localStorage.getItem('token');
  const userStr = localStorage.getItem('user');
  
  if (token && userStr) {
    try {
      authState.token = token;
      authState.user = JSON.parse(userStr);
      authState.isAuthenticated = true;
    } catch (e) {
      console.error('加载用户信息失败:', e);
      clearAuth();
    }
  }
};

// 保存认证信息
const saveAuth = (token, user) => {
  authState.token = token;
  authState.user = user;
  authState.isAuthenticated = true;
  localStorage.setItem('token', token);
  localStorage.setItem('user', JSON.stringify(user));
};

// 清除认证信息
const clearAuth = () => {
  authState.token = null;
  authState.user = null;
  authState.isAuthenticated = false;
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

// 获取后端URL
const getBackendUrl = () => {
  const baseUrl = `${window.location.protocol}//${window.location.host}`;
  return {
    http: baseUrl,
  };
};

// 注册
const register = async (username, password, name) => {
  try {
    const backend = getBackendUrl();
    const response = await fetch(`${backend.http}/api/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password, name }),
    });

    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || '注册失败');
    }

    return data;
  } catch (error) {
    console.error('注册错误:', error);
    throw error;
  }
};

// 登录
const login = async (username, password) => {
  try {
    const backend = getBackendUrl();
    const response = await fetch(`${backend.http}/api/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    
    if (!response.ok) {
      throw new Error(data.error || '登录失败');
    }

    // 保存认证信息
    saveAuth(data.token, data.user);
    
    return data;
  } catch (error) {
    console.error('登录错误:', error);
    throw error;
  }
};

// 登出
const logout = () => {
  clearAuth();
};

// 初始化：加载保存的认证信息
loadAuth();

export default {
  state: authState,
  register,
  login,
  logout,
  clearAuth,
};

