import { AuthMachineContext } from "@/machines/auth-machine"
import { useSelector } from "@xstate/react"

export const LoadingBar = () => {
  const authActor = AuthMachineContext.useActorRef();
  const loaded = useSelector(authActor, ({ context }) => context.loaded);
  return (
    <div className="w-full h-1 fixed top-[53px]">
      <div
        className="bg-yellow-500 h-1"
        style={{
          width: `${loaded}%`,
          transition: loaded > 10 ? "width 3s ease-out" : "width .25 ease-in",
        }}
      />
    </div>
  )
}
