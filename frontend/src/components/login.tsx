import React from "react";
import { useAuth0 } from "@auth0/auth0-react"

export const LoginButton: React.FC = () => {
  const { loginWithRedirect } = useAuth0();
  return ( 
    <button onClick={() => loginWithRedirect()} 
      className="block w-full rounded border border-blue-600 bg-blue-600 px-12 py-3 text-sm font-medium text-white hover:bg-transparent hover:text-white focus:outline-none focus:ring active:text-opacity-75 sm:w-auto"> 
      Log In 
    </button>
  );
};
