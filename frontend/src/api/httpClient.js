const $host = ky.create(
  {
    prefixUrl: BASE_API + "/api",
  }
);

const $authHost = ky.create(
  {
    prefixUrl: BASE_API + "/api",
    hooks: {
      beforeRequest: [
        (req) => {
          const { token } = useUserStore();
          req.headers.set('Authorization', `Bearer ${token}`);
        }
      ]
    }
  },
);

export { $host, $authHost };