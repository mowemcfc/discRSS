export type Feed = {
  feedId: number
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
  userId: number
  username: string
  feedList: Feed[]
  channelList: DiscordChannel[]
}

export default UserAccount
