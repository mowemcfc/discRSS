import { UserAccount, Feed } from '../types/user'
import { FeedRow, NewFeedRow } from './feed-row'
import { useState } from 'react'

export const UserModal = (props: { userAccount: UserAccount }) => {

  const [numRows, setNumRows] = useState(props.userAccount.feedList.length)

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

  const submitForm = (event: React.FormEvent) => {
    event.preventDefault()
    return null
  }

  return (
    <div className="overflow-hidden border rounded-lg">
        <form onSubmit={event => submitForm(event)}>
          <table className="divide-y divide-gray-200">
            <thead key={`FeedTableHeader-${props.userAccount.userID}`} >
              <th
                scope="col"
                className="w-1/12 py-3 text-xs font-bold text-left text-gray-500 uppercase"
              >
                Feed ID
              </th>

              <th
                scope="col"
                className="w-4/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
              >
                Feed Name
              </th>

              <th
                scope="col"
                className="w-3/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
              >
                Feed URL
              </th>

              <th
                scope="col"
                className="w-3/12 px-6 py-3 text-xs font-bold text-left text-gray-500 uppercase"
              >
                Feed Timestamp
              </th>
              <th
                scope="col"
                className="w-1/12 py-3 text-xs font-bold text-left text-gray-500 uppercase"
              >
                &nbsp;
              </th>

            </thead>
            <tbody key={`FeedTableBody-${props.userAccount.userID}`} className="divide-y divide-gray-200">
              {accountRows(props.userAccount.feedList)}
              <NewFeedRow />
            </tbody>
          </table>
        </form>
    </div>
  )
}