function request(url, config) {
  return fetch(url, config)
    .then((res) => {
      if (res.ok) {
        try {
          return res.json();
        } catch (err) {
          return err;
        }
      } else {
        throw res;
      }
    })
    .then((res) => {
      layui.layer.msg(res.message || "success");
      return res;
    })
    .catch((err) => {
      layer.layer.msg(res.message || "请求错误");
    });
}
