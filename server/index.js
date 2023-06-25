const main = () => {
  const wins = document.getElementById('wins')
  const losses = document.getElementById('losses')
  const winRate = document.getElementById('win-rate')
  const lp = document.getElementById('lp')

  const evtSource = new EventSource('/stream');

  evtSource.addEventListener('message', e => {
    const mh = JSON.parse(e.data)
    console.log(mh)

    wins.innerText = mh.wins
    losses.innerText = mh.losses
    winRate.innerText = `${mh.winRate}%`
    lp.innerText = mh.lpGain
  })

  evtSource.addEventListener('error', e => {
    console.log('error', e)
  })
  
  evtSource.addEventListener('open', _ => {
    console.log('The connection has been established.');
  });

  evtSource.addEventListener('error', _ => {
    console.log('An error occurred while attempting to connect.');
  });
}

document.addEventListener('DOMContentLoaded', main)
