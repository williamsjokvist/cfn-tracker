
type LoadingBarProps = {
  percentage: number
}
export const LoadingBar: React.FC<LoadingBarProps> = ( { percentage} ) => {
  return (<div className='w-full h-1 fixed top-[53px]'>
    <div className='bg-yellow-500 h-1' style={{ 
      width: `${percentage}%`,
      transition: (percentage > 10) ? 'width 3s ease-out' : 'width .25 ease-in',
    }}/>
  </div>)
}