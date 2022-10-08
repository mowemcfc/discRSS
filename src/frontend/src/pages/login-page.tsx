import React from "react";
import { LoginButton } from "../components/login";
import { useAuth0 } from "@auth0/auth0-react";
import { Navigate } from "react-router-dom";

export const LoginPage = () => {
  const { isAuthenticated } = useAuth0();
  if (isAuthenticated) {
    return <Navigate replace to="/account" />
  }

  return <LoginButton />;
}