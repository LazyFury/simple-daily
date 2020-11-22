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
      if (res.code == 1) {
        layui.layer.msg(res.message || "success");
        return res;
      }
      throw res;
    })
    .catch((err) => {
      layui.layer.msg(err.message || "请求错误");
      throw err;
    });
}

function reload() {
  setTimeout(() => {
    location.reload();
  }, 1000);
}

function parentReload() {
  setTimeout(() => {
    window.parent.location.reload();
  }, 1000);
}

function reloadIframe(api) {
  return api()
    .then(() => {
      parentReload();
    })
    .catch(() => {
      reload();
    });
}
