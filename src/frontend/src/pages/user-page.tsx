import React, { useState, useEffect } from "react";
import { UserProfile } from "../components/profile";
import { useAuth0 } from "@auth0/auth0-react";
import { Navigate } from "react-router-dom";
import { UserModal } from "../components/user-modal";
import { SiteBanner } from '../components/site-banner'
import UserAccount from '../types/user'

export const UserPage = () => {
  const {
    isAuthenticated,
    isLoading,
    user,
    getAccessTokenSilently,
  } = useAuth0();

  const [account, setAccount] = useState<UserAccount>({ userID: -1, username: '', feedList: [], channelList: [] })

  
  useEffect(() => {
    const url = process.env.REACT_APP_APIGW_ENDPOINT!
    const fetchUser = async (id: number) => {
      const accessToken = await getAccessTokenSilently()
      const user = await fetch(`${url}user?userID=${id}`, {
        headers: {
          authorization: `Bearer ${accessToken}`
        },
      })
        .then(res => { return res.json() })
        .then(data => { return JSON.parse(data["body"]) })

      setAccount(user)
    }

    fetchUser(1)
  }, [user?.email, getAccessTokenSilently])

  if (isLoading) {
    return <div> Loading your user profile... </div>
  }

  if(!isAuthenticated) {
    return <Navigate replace to="/login" />
  }

  return (
    <div>
      <div className="flex justify-center px-2 py-3 bg-blue-300"> 
        <SiteBanner />
      </div>
      <div className="bg-slate-200 mx-auto px-4 py-32 lg:items-center h-screen">
        <UserProfile />
        <UserModal/>
      </div>
    </div>
  )
}
