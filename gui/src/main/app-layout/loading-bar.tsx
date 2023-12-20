import { AuthMachineContext } from "@/machines/auth-machine"
import { useSelector } from "@xstate/react"

export const LoadingBar = () => {
  const authActor = AuthMachineContext.useActorRef();
  const progress = useSelector(authActor, ({ context }) => context.progress);
  return (
    <div className="w-full h-1 fixed top-[53px]">
      <div
        className="bg-yellow-500 h-1"
        style={{
          width: `${progress}%`,
          transition: progress > 10 ? "width 3s ease-out" : "width .25 ease-in",
        }}
      />
    </div>
  )
}
