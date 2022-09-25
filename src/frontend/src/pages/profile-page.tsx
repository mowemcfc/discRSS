import React from "react";
import { UserProfile } from "../components/profile";
import { useAuth0 } from "@auth0/auth0-react";
import { Navigate } from "react-router-dom";
import { UserAccount } from "../components/user-account";

export const ProfilePage = () => {
  const {
    isAuthenticated,
    isLoading,
  } = useAuth0();

  if (isLoading) {
    return <div> Loading your user profile... </div>
  }

  if(!isAuthenticated) {
    return <Navigate replace to="/login" />
  }

  return (
    <div>
      <UserProfile />
      <UserAccount id={1} />
    </div>
  )
}
