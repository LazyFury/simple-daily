function request(url, config) {
  return fetch(url, config).then((res) => {
    if (res.ok) {
      try {
        return res.json();
      } catch (err) {
        return err;
      }
    } else {
      throw res;
    }
  });
}
