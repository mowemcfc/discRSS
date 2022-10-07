
import { UserAccount, Feed } from '../types/user'
import { FeedRow } from './feed-row'
import React from 'react'

export const UserModal = (props: { userAccount: UserAccount }) => {

  if (!props.userAccount) {
    return (
      <div>
        <div>Loading your user account ...</div>
      </div>
    )
  }

  const accountRows  = (feedList: Feed[]) => feedList.map((feed: Feed) => {
    return (
      <FeedRow feed={feed} />
    )
  })

  return (
    <div>
        <table>
          <tbody>
            {accountRows(props.userAccount.feedList)}
          </tbody>
        </table>
    </div>
  )
}