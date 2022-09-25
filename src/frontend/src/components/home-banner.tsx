import React from "react";
import { Link } from "react-router-dom"

export const HomeBanner = () => {

  return (
    <div> 
      Hello!
      <Link to="/login"> Log In Page </Link>
    </div>
  );
};
