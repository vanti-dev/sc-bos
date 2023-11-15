const api = import.meta.env.VITE_AUTH_URL || '';

const localLogin = (username, password) => {
  const formData = new FormData();
  formData.set('grant_type', 'password');
  formData.set('username', username);
  formData.set('password', password);
  return fetch(`${api}/oauth2/token`, {
    method: 'POST',
    body: formData
  });
};

export default localLogin;
