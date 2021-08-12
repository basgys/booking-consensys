import { useContext } from "react";
import { AuthContext } from "src/user/auth";
import useSWR from "swr";

const Rooms = () => {
  const auth = useContext(AuthContext)
  const { data, error } = useSWR('/booking/rooms', auth.fetcher)

  if (error) return <div>{error.message}</div>
  if (!data) return <div>loading...</div>

  return (
    <ul className="rooms">
      {data.rooms.map((room: any) => (
        <li key={room.ref}>
          <a href={`/rooms/${room.ref}`}>{room.ref}</a>
        </li>
      ))}
    </ul>
  );
}

export default Rooms;
