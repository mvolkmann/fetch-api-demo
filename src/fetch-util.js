const USE_HTTPS = false;

const protocol = USE_HTTPS ? 'https' : 'http';
const isLocal = window.location.href.startsWith('http://localhost:');

// Change this to match the URL prefix of your REST services.
// If your project uses REST services with more than one URL prefix,
// drop the use URL_PREFIX and just pass full URLs into the functions.
const URL_PREFIX = isLocal ? `${protocol}://localhost:3000/` : '/';

// If there are any common options that are
// desired for all REST calls, place them here.
const options = {};

const headers = {'Content-Type': 'application/json'};

// Can't name this "delete" because that is a JavaScript keyword.
export function deleteResource(urlSuffix) {
  const url = URL_PREFIX + urlSuffix;
  return fetch(url, {...options, method: 'DELETE'});
}

export async function getJson(urlSuffix) {
  const url = URL_PREFIX + urlSuffix;
  const res = await fetch(url, options);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function getText(urlSuffix) {
  const url = URL_PREFIX + urlSuffix;
  const res = await fetch(url, options);
  if (!res.ok) throw new Error(await res.text());
  return res.text();
}

export function postJson(urlSuffix, obj) {
  const url = URL_PREFIX + urlSuffix;
  const body = JSON.stringify(obj);
  return fetch(url, {...options, method: 'POST', headers, body});
}

export function putJson(urlSuffix, obj) {
  const url = URL_PREFIX + urlSuffix;
  const body = JSON.stringify(obj);
  return fetch(url, {...options, method: 'PUT', headers, body});
}
