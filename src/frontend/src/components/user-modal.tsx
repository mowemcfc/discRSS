
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
      <FeedRow key={`FeedRow-${feed.feedID}`} feed={feed} />
    )
  })

  return (
    <div className="overflow-hidden border rounded-lg">
        <table className="min-w-full divide-y divide-gray-200">
          <thead key={`FeedTableHeader-${props.userAccount.userID}`} >
            <th
              scope="col"
              className="px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
            >
              FeedID
            </th>

            <th
              scope="col"
              className="px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
            >
              Feed Name
            </th>

            <th
              scope="col"
              className="px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
            >
              Feed URL
            </th>

            <th
              scope="col"
              className="px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
            >
              Feed Timestamp
            </th>
          </thead>
          <tbody key={`FeedTableBody-${props.userAccount.userID}`} className="divide-y divide-gray-200">
            {accountRows(props.userAccount.feedList)}
          </tbody>
        </table>
    </div>
  )
}