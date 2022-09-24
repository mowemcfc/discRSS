import React from "react";
import { LoginButton } from "../components/login";
import { useAuth0, User } from "@auth0/auth0-react";
import { Navigate, redirect } from "react-router-dom";

export const LoginPageLoader = () => {
  const { user } = useAuth0();
  return user
}

export const LoginPageAction = ({ user }: User) => {
  if (user) {
    return redirect("/");
  }
}

export const LoginPage = () => {
  const { user } = useAuth0();
  if (user) {
    return <Navigate replace to="/" />
  }

  return <LoginButton />;
}