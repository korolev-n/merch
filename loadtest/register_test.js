import http from 'k6/http';
import { check } from 'k6';
import { SharedArray } from 'k6/data';

export const options = {
  vus: 100, // 100000 пользователей
  duration: '1s',
};

const users = new SharedArray('users', () => {
  return Array.from({ length: 100 }, (_, i) => ({
    username: `user${i}@loadtest.io`, 
    password: 'SuperSecure123!', 
  }));
});

export default function () {
  const userIndex = __VU - 1; 

  const user = users[userIndex]; 

  const res = http.post(
    'http://localhost:8080/api/auth',
    JSON.stringify({ username: user.username, password: user.password }),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { phase: 'register' },
    }
  );

  check(res, {
    'status is 200': (r) => r.status === 200,
    'token exists': (r) => !!r.json('token'),
  });
}
