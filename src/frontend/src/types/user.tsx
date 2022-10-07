export type Feed = {
  feedID: number
  title: string
  url: string
  timeFormat: string
}

export type DiscordChannel = {
  channelName: string
  serverName: string
  channelID: number
}

export type UserAccount = {
  userID: number
  username: string
  feedList: Feed[]
  channelList: DiscordChannel[]
}

export default UserAccount