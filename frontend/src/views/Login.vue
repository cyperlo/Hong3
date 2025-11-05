<template>
  <div class="login-container">
    <!-- 主内容区域 -->
    <div class="login-content">
      <div class="login-wrapper">
        <h1 class="app-title">红3</h1>
        <p class="app-subtitle">欢迎登录</p>
        
        <div class="button-area">
          <button
            @click="handleLoginClick"
            class="action-btn login-btn"
          >
            登录
          </button>
          <button
            @click="handleRegisterClick"
            class="action-btn register-btn"
          >
            注册
          </button>
        </div>
      </div>
    </div>

    <!-- 登录弹窗 -->
    <div v-if="showLogin" class="modal-overlay" @click.self="showLogin = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>登录</h2>
          <button class="close-btn" @click="showLogin = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <input
              v-model="loginForm.username"
              type="text"
              placeholder="账号"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              class="form-input"
              @keyup.enter="handleLogin"
            />
          </div>
          <div v-if="loginError" class="error-message">
            {{ loginError }}
          </div>
          <button
            @click="handleLogin"
            :disabled="loggingIn || !loginForm.username || !loginForm.password"
            class="btn btn-primary"
          >
            {{ loggingIn ? '登录中...' : '登录' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 注册弹窗 -->
    <div v-if="showRegister && !showLogin" class="modal-overlay" @click.self="showRegister = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>注册</h2>
          <button class="close-btn" @click="showRegister = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <input
              v-model="registerForm.username"
              type="text"
              placeholder="账号"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <input
              v-model="registerForm.password"
              type="password"
              placeholder="密码"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <input
              v-model="registerForm.name"
              type="text"
              placeholder="昵称（游戏中的显示名称）"
              class="form-input"
              @keyup.enter="handleRegister"
            />
          </div>
          <div v-if="registerError" class="error-message">
            {{ registerError }}
          </div>
          <button
            @click="handleRegister"
            :disabled="registering || !registerForm.username || !registerForm.password || !registerForm.name"
            class="btn btn-primary"
          >
            {{ registering ? '注册中...' : '注册' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import authStore from '../store/authStore';

const emit = defineEmits(['login-success']);
const showLogin = ref(false);
const showRegister = ref(false);
const loggingIn = ref(false);
const registering = ref(false);
const loginError = ref('');
const registerError = ref('');

const loginForm = reactive({
  username: '',
  password: '',
});

const registerForm = reactive({
  username: '',
  password: '',
  name: '',
});


// 处理登录按钮点击
const handleLoginClick = (e) => {
  e.preventDefault();
  e.stopPropagation();
  console.log('Login button clicked, showLogin before:', showLogin.value);
  showLogin.value = true;
  showRegister.value = false;
  console.log('Login button clicked, showLogin after:', showLogin.value);
};

// 处理注册按钮点击
const handleRegisterClick = (e) => {
  e.preventDefault();
  e.stopPropagation();
  console.log('Register button clicked, showRegister before:', showRegister.value);
  showRegister.value = true;
  showLogin.value = false;
  registerError.value = '';
  console.log('Register button clicked, showRegister after:', showRegister.value);
};

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    loginError.value = '请填写账号和密码';
    return;
  }

  loginError.value = '';
  loggingIn.value = true;

  try {
    await authStore.login(loginForm.username, loginForm.password);
    showLogin.value = false;
    loginForm.username = '';
    loginForm.password = '';
    loginError.value = '';
    emit('login-success');
  } catch (err) {
    loginError.value = err.message || '登录失败';
  } finally {
    loggingIn.value = false;
  }
};

const handleRegister = async () => {
  if (!registerForm.username || !registerForm.password || !registerForm.name) {
    registerError.value = '请填写所有信息';
    return;
  }

  registerError.value = '';
  registering.value = true;

  try {
    // 注册
    await authStore.register(registerForm.username, registerForm.password, registerForm.name);
    
    // 注册成功后，等待一小段时间确保后端处理完成，然后自动登录
    await new Promise(resolve => setTimeout(resolve, 300));
    
    // 自动登录
    await authStore.login(registerForm.username, registerForm.password);
    
    // 关闭弹窗并清空表单
    showRegister.value = false;
    showLogin.value = false;
    registerForm.username = '';
    registerForm.password = '';
    registerForm.name = '';
    registerError.value = '';
    
    // 触发登录成功事件
    emit('login-success');
  } catch (err) {
    // 如果注册失败，显示注册错误
    // 如果注册成功但登录失败，显示登录错误
    if (err.message && err.message.includes('登录')) {
      registerError.value = '注册成功，但登录失败: ' + (err.message || '未知错误');
    } else {
      registerError.value = err.message || '注册失败';
    }
  } finally {
    registering.value = false;
  }
};

</script>

<style scoped>
.login-container {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  min-width: 100vw !important;
  min-height: 100vh !important;
  max-width: 100vw !important;
  max-height: 100vh !important;
  overflow: hidden !important;
  z-index: 99999 !important;
  box-sizing: border-box !important;
  margin: 0 !important;
  padding: 0 !important;
}


/* 主内容区域 */
.login-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  box-sizing: border-box;
}

.login-wrapper {
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.app-title {
  font-size: 3rem;
  font-weight: bold;
  color: #d32f2f;
  margin: 0 0 0.5rem 0;
  letter-spacing: 2px;
}

.app-subtitle {
  font-size: 1.1rem;
  color: #666;
  margin: 0 0 3rem 0;
}

.button-area {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
}

.action-btn {
  width: 100%;
  padding: 0.875rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  touch-action: manipulation;
}

.login-btn {
  background-color: #1b5e20;
  color: white;
}

.login-btn:hover {
  background-color: #2e7d32;
}

.login-btn:active {
  background-color: #1b5e20;
  transform: scale(0.98);
}

.register-btn {
  background-color: #d32f2f;
  color: white;
}

.register-btn:hover {
  background-color: #f44336;
}

.register-btn:active {
  background-color: #d32f2f;
  transform: scale(0.98);
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 450px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  color: #333;
  font-weight: bold;
}

.close-btn {
  background: none;
  border: none;
  font-size: 2rem;
  color: #666;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.2rem;
}

.form-input {
  width: 100%;
  padding: 0.875rem 1rem;
  font-size: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  transition: all 0.2s;
  box-sizing: border-box;
  background: white;
  color: #333;
}

.form-input::placeholder {
  color: #999;
}

.form-input:focus {
  outline: none;
  border-color: #1b5e20;
  box-shadow: 0 0 0 3px rgba(27, 94, 32, 0.1);
}

.btn {
  width: 100%;
  padding: 1rem;
  font-size: 1.1rem;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 0.5rem;
  box-sizing: border-box;
  display: block;
  position: relative;
  z-index: 1;
}

.btn-primary {
  background-color: #1b5e20;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2e7d32;
}

.btn-primary:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-primary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
  transform: none;
  color: #999;
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: #ffebee;
  color: #c62828;
  border-radius: 8px;
  border-left: 4px solid #c62828;
  font-size: 0.9rem;
}

/* 移动端优化 */
@media screen and (max-width: 768px) {
  .login-content {
    padding: 1rem;
  }
  
  .app-title {
    font-size: 2.5rem;
  }
  
  .app-subtitle {
    font-size: 1rem;
    margin-bottom: 2rem;
  }
  
  .modal-content {
    max-width: 90%;
  }
}

</style>
