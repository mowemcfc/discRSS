import React, { useState, useEffect } from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { UserProfile } from "../components/profile";
import { Navigate } from "react-router-dom";
import { Feed } from '../types/user'
import { FeedList } from "../components/feed-list";
import { SiteBanner } from '../components/site-banner'
import { handleErrors } from "../utils";
import UserAccount from '../types/user'

export const UserPage: React.FC = () => {
  const {
    isAuthenticated,
    getAccessTokenSilently,
    getIdTokenClaims,
  } = useAuth0();


  const [userFeedList, setUserFeedList] = useState<Feed[]>([])
  const [isFirstLogin, setIsFirstLogin] = useState<boolean | null>(null);

  const fetchUser = async (id: number) => {
    const accessToken = await getAccessTokenSilently()
    const resp = await fetch(`${process.env.REACT_APP_APIGW_ENDPOINT!}user/${id}/feeds`, {
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(handleErrors)
      .then(res => res.json())
      .then(data => data.Body)
    setUserFeedList(resp)
  }

  const checkFirstLogin = async () => {
    if(isAuthenticated) {
      const claims = await getIdTokenClaims()
      console.log(claims)
      setIsFirstLogin(false) // TODO: properly check this
    }
  }


  useEffect(() => {
    checkFirstLogin()
    fetchUser(10) // TODO: get this value dynamically from auth0 ID
  }, [])

  if(!isAuthenticated) {
    return <Navigate replace to="/login" />
  }

  if(isFirstLogin) {
    return <Navigate replace to="/register"/>
  }

  if (!userFeedList) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex-initial justify-center"> 
        <SiteBanner />
      </div>
      <div className="bg-gray-900 px-4 py-32 h-full min-h-screen">
        <FeedList feedList={Object.values(userFeedList)} userId={10}/>
      </div>
    </div>
  )
}
