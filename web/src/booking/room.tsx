import { useContext, useState } from "react";
import { useParams } from "react-router-dom";
import { parseURL } from "src/lib/url";
import { AuthContext } from "src/user/auth";
import useSWR from "swr";

interface RoomParams {
  id: string
}

interface State {
  from: Date
  to: Date
}

const Room = () => {
  const { id } = useParams<RoomParams>();
  const auth = useContext(AuthContext)
  const now = new Date(Date.now())
  const to = new Date(now.valueOf())
  to.setDate(to.getDate() + 7); // 7 days
  const [state] = useState<State>({
    from: now,
    to: to,
  })

  const url = parseURL(`/booking/rooms/${id}/availabilities`)
  url.searchParams.set('from', state.from.toISOString())
  url.searchParams.set('to', state.to.toISOString())
  const { data, error } = useSWR(url.toString(), auth.fetcher)

  if (error) return <div>{error.message}</div>
  if (!data) return <div>loading...</div>

  return (
    <div className="room">
      <h1>Room #{id}</h1>

      <h2>Availabilities</h2>
      <ul className="room__availabilities">
        {data.availabilities?.map((a: any, i: number) => {
          const from = new Date(a.from)
          const to = new Date(a.to)
          return (
            <li key={i}>
              <time dateTime={from.toISOString()}>
                {from.toLocaleString()}
              </time> –
              <time dateTime={to.toISOString()}>
                {to.toLocaleString()}
              </time>
            </li>
          )
        })}
      </ul>
    </div>
  );
}

export default Room;
