import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export const options = {
  scenarios: {
    constant_rps: {
      executor: 'constant-arrival-rate',
      rate: 100, // 1000 RPS по всей системе
      timeUnit: '1s',
      duration: '10s',
      preAllocatedVUs: 100, // 1000
      maxVUs: 200, // 2000
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<50'],
    http_req_failed: ['rate<0.0001'],
  },
};

// на 100 пользователей
const users = new SharedArray('users', () => {
  return Array.from({ length: 100 }, (_, i) => ({
    username: `user_${i}`,
    password: 'Password123',
  }));
});

function register(user) {
  const res = http.post(
    'http://localhost:8080/api/auth',
    JSON.stringify({ username: user.username, password: user.password }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(res, {
    'auth status 200': (r) => r.status === 200,
    'token returned': (r) => r.json('token') !== undefined,
  });
  return res.json('token');
}

function sendCoin(token) {
  const payload = JSON.stringify({
    toUser: 'user_0', // отправим всем coin пользователю #0
    amount: 1,
  });
  const res = http.post('http://localhost:8080/api/sendCoin', payload, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });
  check(res, {
    'sendCoin status 200': (r) => r.status === 200 || r.status === 400 || r.status === 404,
  });
}

function buyItem(token) {
  const items = ['cup', 'pen', 'book', 't-shirt'];
  const item = items[Math.floor(Math.random() * items.length)];
  const res = http.get(`http://localhost:8080/api/buy/${item}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  check(res, {
    'buyItem status 200': (r) => r.status === 200 || r.status === 400 || r.status === 404,
  });
}

function getInfo(token) {
  const res = http.get('http://localhost:8080/api/info', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  check(res, {
    'getInfo status 200': (r) => r.status === 200,
    'has coins': (r) => r.json('coins') !== undefined,
  });
}

export default function () {
  const user = users[__VU % users.length];
  const token = register(user);

  if (token) {
    sendCoin(token);
    buyItem(token);
    getInfo(token);
  }

  sleep(1); // пауза между итерациями (важно при реальной нагрузке)
}
