import { useLocation } from "react-router-dom";

function AuthorizationPage() {
  const { hash } = useLocation();
  return (
    <>
      {hash == "auth" ? (
        <AuthForm />
      ):(
        <RegForm />
      )}
    </>

  );
}

export default AuthorizationPage;