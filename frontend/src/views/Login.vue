<template>
  <div class="login-container">
    <div class="login-content">
      <div class="logo-section">
        <h1 class="app-title">
          <span class="title-red">红</span><span class="title-number">3</span>
        </h1>
        <p class="subtitle">欢迎来到红3游戏</p>
      </div>

      <div class="login-form-container">
        <div class="login-form">
          <h2>登录</h2>
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
          <button
            @click="handleLogin"
            :disabled="loggingIn || !loginForm.username || !loginForm.password"
            class="btn btn-primary"
          >
            {{ loggingIn ? '登录中...' : '登录' }}
          </button>
          <div class="form-footer">
            <span>还没有账号？</span>
            <button class="link-btn" @click="showRegister = true; registerError = ''">立即注册</button>
          </div>
        </div>

        <div v-if="loginError" class="error-message">
          {{ loginError }}
        </div>
      </div>
    </div>

    <!-- 注册弹窗 -->
    <div v-if="showRegister" class="modal-overlay" @click.self="showRegister = false">
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

const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    loginError.value = '请填写账号和密码';
    return;
  }

  loginError.value = '';
  loggingIn.value = true;

  try {
    await authStore.login(loginForm.username, loginForm.password);
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
    await authStore.register(registerForm.username, registerForm.password, registerForm.name);
    // 注册成功后自动登录
    try {
      await authStore.login(registerForm.username, registerForm.password);
      showRegister.value = false; // 关闭弹窗
      // 清空注册表单
      registerForm.username = '';
      registerForm.password = '';
      registerForm.name = '';
      registerError.value = '';
      emit('login-success');
    } catch (err) {
      registerError.value = '注册成功，但登录失败: ' + (err.message || '未知错误');
    }
  } catch (err) {
    registerError.value = err.message || '注册失败';
  } finally {
    registering.value = false;
  }
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.login-content {
  width: 100%;
  max-width: 420px;
  background: white;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.logo-section {
  background: linear-gradient(135deg, #d32f2f 0%, #b71c1c 100%);
  color: white;
  padding: 2rem;
  text-align: center;
}

.app-title {
  font-size: 3rem;
  font-weight: bold;
  margin: 0;
  letter-spacing: 2px;
}

.title-red {
  color: #ffeb3b;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.title-number {
  color: white;
  margin-left: 0.2rem;
}

.subtitle {
  margin-top: 0.5rem;
  font-size: 1rem;
  opacity: 0.9;
}

.login-form-container {
  padding: 2rem;
}

.login-form h2,
.register-form h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.8rem;
  color: #333;
  text-align: center;
}

.form-group {
  margin-bottom: 1.2rem;
}

.form-input {
  width: 100%;
  padding: 1rem;
  font-size: 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  transition: border-color 0.3s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #d32f2f;
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
}

.btn-primary {
  background: linear-gradient(135deg, #d32f2f 0%, #b71c1c 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(211, 47, 47, 0.4);
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.form-footer {
  margin-top: 1.5rem;
  text-align: center;
  color: #666;
  font-size: 0.9rem;
}

.link-btn {
  background: none;
  border: none;
  color: #d32f2f;
  cursor: pointer;
  font-weight: 600;
  margin-left: 0.5rem;
  padding: 0;
  text-decoration: underline;
}

.link-btn:hover {
  color: #b71c1c;
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

/* 横屏手机端优化 */
@media (max-width: 768px) and (orientation: landscape) {
  .login-container {
    padding: 0.5rem;
  }

  .login-content {
    max-width: 100%;
    border-radius: 12px;
  }

  .logo-section {
    padding: 1.5rem;
  }

  .app-title {
    font-size: 2.5rem;
  }

  .login-form-container {
    padding: 1.5rem;
  }

  .login-form h2,
  .register-form h2 {
    font-size: 1.5rem;
  }
}

/* 竖屏手机端 */
@media (max-width: 768px) and (orientation: portrait) {
  .login-container {
    padding: 1rem;
  }

  .app-title {
    font-size: 2.5rem;
  }
}

/* 注册弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: modalSlideIn 0.3s ease-out;
}

@keyframes modalSlideIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
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
  transition: all 0.3s;
}

.close-btn:hover {
  background: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 1.5rem;
}

.modal-body .form-group {
  margin-bottom: 1rem;
}

.modal-body .btn {
  margin-top: 1rem;
}
</style>

