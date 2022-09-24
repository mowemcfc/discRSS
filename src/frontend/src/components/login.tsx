import React from "react";
import { useAuth0 } from "@auth0/auth0-react"
import { redirect } from "react-router-dom";

export const LoginButton = () => {
  const { loginWithRedirect } = useAuth0();

  return <button onClick={() => loginWithRedirect()}> Log In </button>;
};

export const LoginPage = () => {
  return <LoginButton />;
}
