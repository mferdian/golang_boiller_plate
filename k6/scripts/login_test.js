import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,          
  duration: '30s', 
};

export default function () {
  const url = 'http://localhost:8000/api/login';

  const payload = JSON.stringify({
    email: 'michael01@gmail.com',
    password: 'password123',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'has access token': (r) => {
      const body = r.json();
      return body?.data?.access_token !== undefined;
    },
  });

  sleep(1);
}
