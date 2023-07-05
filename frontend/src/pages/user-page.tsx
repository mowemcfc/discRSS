import React, { useState, useEffect } from "react";
import { useAuth0 } from "@auth0/auth0-react";
import { UserProfile } from "../components/profile";
import { Navigate } from "react-router-dom";
import UserAccount, { Feed } from '../types/user'
import { FeedList } from "../components/feed-list";
import { SiteBanner } from '../components/site-banner'
import { handleErrors } from "../utils";

export const UserPage: React.FC = () => {
  const {
    isAuthenticated,
    getAccessTokenSilently,
    getIdTokenClaims,
  } = useAuth0();


  const [userFeedList, setUserFeedList] = useState<Feed[]>([])
  const [isFirstLogin, setIsFirstLogin] = useState<boolean | null>(null)
  const [userId, setUserId] = useState<number>(0)

  const getUser = async (id: number): Promise<UserAccount> => {
    const accessToken = await getAccessTokenSilently()
    const user = await fetch(`${process.env.REACT_APP_APIGW_ENDPOINT!}user/${id}`, {
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(res => res.json())
      .then(data => data.Body as UserAccount)
    return user
  }

  const createUser = async (id: number): Promise<UserAccount> => {
    const accessToken = await getAccessTokenSilently()
    const user = await fetch(`${process.env.REACT_APP_APIGW_ENDPOINT!}user`, {
      method: 'POST',
      body: JSON.stringify({
        userId: id.toString(),
        username: 'mower'
      }),
      headers: {
        authorization: `Bearer ${accessToken}`
      },
    })
      .then(handleErrors)
      .then(res => res.json())
      .then(data => data.Body as UserAccount)
    return user
  }

  const getOrCreateUser = async(id: number): Promise<UserAccount> => {
    var user: UserAccount 
    user = await getUser(id)
    const isEmpty = Object.values(user).every(x => x === null || x === '');
    return isEmpty ? createUser(id) : user
  }

  const checkFirstLogin = async () => {
    if(isAuthenticated) {
      const claims = await getIdTokenClaims()
      // Auth0 JWT sub is in format oauth2|<provider>|<id>.
      // We will only ever use a single oauth provider, so the ID is solely unique.
      const id = claims!["sub"].split('|')[2]
      const user = await getOrCreateUser(id)
      setUserId(id)
      setUserFeedList(Object.values(user.feedList))
    }
  }


  useEffect(() => {
    checkFirstLogin()
    getUser(userId) // TODO: get this value dynamically from auth0 ID
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
      <div className="flex flex-row justify-between bg-gray-900 grid-cols-2 flex-nowrap py-3"> 
        <SiteBanner />
        <UserProfile />
      </div>
      <div className="bg-gray-900 flex justify-center px-4 py-8 h-full min-h-screen">
        <FeedList feedList={Object.values(userFeedList)} userId={userId}/>
      </div>
    </div>
  )
}
