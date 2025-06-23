export const routes = {
  public: {
    main: "/",
    auth: {
      route: "/auth",
      hashes: {
        auth: "#login",
        reg: "#registration",
      }
    },
  },
  user: {
    home: "/home",
    album: "/album/:id",
    album_edit: "/album/:id/edit",
    artist: "/artist/:id",
    playlist: "/playlist/:id",
  }
}
