export type Feed = {
  feedId: number
  title: string
  url: string
  timeFormat: string
}

export type NewFeedParams = {
  title: string
  url: string
}

export type DiscordChannel = {
  channelName: string
  serverName: string
  channelID: number
}

export type UserAccount = {
  userId: number
  username: string
  feedList: { [id:string]: Feed }
  channelList: { [id:string]: DiscordChannel }
}

export default UserAccount
