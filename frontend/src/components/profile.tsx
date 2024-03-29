import React from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { LogoutButton } from "./logout";

export const UserProfile = () => {
  const {
    user,
    isLoading,
  } = useAuth0();

  if(isLoading) {
    return <div> Loading your user information </div>
  }

  return (
    <div className="basis-1/4 flex place-content-end text-white px-10">
      <div className="px-5">
        <h2>{user?.name}</h2>
        <LogoutButton />
      </div>
      <img className="w-16 h-16" src={user?.picture} alt={user?.name} />
    </div>
  );
};
