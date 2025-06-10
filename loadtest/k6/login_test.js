import http from 'k6/http';
import { check } from 'k6';
import { SharedArray } from 'k6/data';

export const options = {
  scenarios: {
    constant_rps: {
      executor: 'constant-arrival-rate',
      rate: 100, // надо 1000 RPS
      timeUnit: '1s',
      duration: '10s',
      preAllocatedVUs: 100, // надо 1000 при 1000 пользователей 
      maxVUs: 200, // надо 2000 при 1000 пользователей 
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<50'], // 95-й перцентиль < 50 мс
    http_req_failed: ['rate<0.0001'], // < 0.01% ошибок
  },
};

// Используем SharedArray для хранения учетных данных пользователей (с расчетом на 100к)
const users = new SharedArray('users', () => {
  return Array.from({ length: 100 }, (_, i) => ({
    username: `user_${i}`,
    password: 'Password123',
  }));
});

export default function () {
  // Используем __VU для выбора уникального пользователя
  const user = users[__VU % users.length];

  const res = http.post(
    'http://localhost:8080/api/auth',
    JSON.stringify({ username: user.username, password: user.password }),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );

  check(res, {
    'status is 200': (r) => r.status === 200,
    'valid response': (r) => {
      const data = r.json();
      return data && typeof data.token === 'string';
    },
  });
}
