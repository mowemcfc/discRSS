import React, { useState, useEffect } from "react";
import { UserProfile } from "../components/profile";
import { useAuth0 } from "@auth0/auth0-react";
import { Navigate } from "react-router-dom";
import { FeedList } from "../components/feed-list";
import { SiteBanner } from '../components/site-banner'
import { handleErrors } from "../utils";
import UserAccount from '../types/user'

export const UserPage: React.FC = () => {
  const {
    isAuthenticated,
    getAccessTokenSilently,
  } = useAuth0();

  const [userData, setUser] = useState<UserAccount>({ userId: -1, username: '', feedList: [], channelList: [] })

  useEffect(() => {
    const fetchUser = async (id: number) => {
      const accessToken = await getAccessTokenSilently()
      const resp = await fetch(`${process.env.REACT_APP_APIGW_ENDPOINT!}user?userId=${id}`, {
        headers: {
          authorization: `Bearer ${accessToken}`
        },
      })
        .then(handleErrors)
        .then(res => res.json())
        .then(data => data.Body)
      

      setUser(resp)
    }

    fetchUser(1) // TODO: get this value dynamically from auth0 ID 
  }, [])

  if(!isAuthenticated) {
    return <Navigate replace to="/login" />
  }

  if (!userData) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex-initial justify-center px-2 py-3 bg-blue-300"> 
        <SiteBanner />
      </div>
      <div className="flex-1 bg-slate-200 mx-auto px-4 py-32 lg:items-center h-full min-h-screen">
        <FeedList feedList={userData.feedList} userId={userData.userId}/>
      </div>
    </div>
  )
}
