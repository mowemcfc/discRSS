
import UserAccount from '../types/user'
import React from 'react'

export const UserModal = (props: { userAccount: UserAccount }) => {

  if (!props.userAccount) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  return (
    <div>
      <h1>User account for {props.userAccount.username} with user ID {props.userAccount.userID}</h1>
    </div>
  )
}