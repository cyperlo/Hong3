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
    const url = `${backend.http}/api/register`;
    console.log('注册请求 URL:', url);
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password, name }),
    });

    console.log('注册响应状态:', response.status, response.statusText);
    console.log('注册响应 URL:', response.url);

    // 先读取响应文本（Response body只能读取一次）
    const text = await response.text();
    console.log('注册响应内容:', text);
    
    // 检查响应状态
    if (!response.ok) {
      // 尝试解析错误信息
      let errorMessage = '注册失败';
      if (text && text.trim() !== '') {
        try {
          const errorData = JSON.parse(text);
          errorMessage = errorData.error || errorMessage;
        } catch (e) {
          errorMessage = text || response.statusText || `HTTP ${response.status}`;
        }
      } else {
        errorMessage = response.statusText || `HTTP ${response.status}`;
      }
      throw new Error(errorMessage);
    }

    // 检查响应是否为空
    if (!text || text.trim() === '') {
      console.error('注册响应为空');
      throw new Error('服务器返回空响应');
    }

    // 检查响应内容类型
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      console.error('注册响应不是JSON:', contentType, '响应内容:', text);
      throw new Error('服务器返回格式错误');
    }

    // 解析JSON
    let data;
    try {
      data = JSON.parse(text);
    } catch (e) {
      console.error('JSON解析失败:', e, '响应内容:', text);
      throw new Error('服务器返回数据格式错误');
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
    const url = `${backend.http}/api/login`;
    console.log('登录请求 URL:', url);
    
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    console.log('登录响应状态:', response.status, response.statusText);
    console.log('登录响应 URL:', response.url);

    // 先读取响应文本（Response body只能读取一次）
    const text = await response.text();
    console.log('登录响应内容:', text);
    
    // 检查响应状态
    if (!response.ok) {
      // 尝试解析错误信息
      let errorMessage = '登录失败';
      if (text && text.trim() !== '') {
        try {
          const errorData = JSON.parse(text);
          errorMessage = errorData.error || errorMessage;
        } catch (e) {
          // 如果不是JSON，使用原始文本或状态文本
          errorMessage = text || response.statusText || `HTTP ${response.status}`;
        }
      } else {
        errorMessage = response.statusText || `HTTP ${response.status}`;
      }
      throw new Error(errorMessage);
    }

    // 检查响应是否为空
    if (!text || text.trim() === '') {
      console.error('登录响应为空');
      throw new Error('服务器返回空响应');
    }

    // 检查响应内容类型
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      console.error('登录响应不是JSON:', contentType, '响应内容:', text);
      throw new Error('服务器返回格式错误');
    }

    // 解析JSON
    let data;
    try {
      data = JSON.parse(text);
    } catch (e) {
      console.error('JSON解析失败:', e, '响应内容:', text);
      throw new Error('服务器返回数据格式错误');
    }
    
    // 验证必要字段
    if (!data.token || !data.user) {
      throw new Error('服务器返回数据不完整');
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

