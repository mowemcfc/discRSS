import React from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { LogoutButton } from "./logout";

export const UserProfile = () => {
  const {
    user,
    isAuthenticated,
    isLoading,
    getAccessTokenSilently,
  } = useAuth0();

  if(isLoading) {
    return <div> Loading your user information </div>
  }

  return (
    <div>
      <img src={user?.picture} alt={user?.name} />
      <h2>{user?.name}</h2>
      <p>{user?.email}</p>
      <LogoutButton />
    </div>
  );
};
