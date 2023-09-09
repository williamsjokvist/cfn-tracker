class StatList extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow( { mode: "open" } )
  }

  addItem(title, value) {
    const item = document.createElement('stat-list-item')
    item.title = title
    item.value = value
    item.part = 'list-item'

    this.shadow.appendChild(item)
  }

  removeItem(title) {
    const item = this.shadow.querySelector(`stat-list-item[title="${title}"]`)
    if (item != null) {
      item.remove()
    }
  }

  updateItem(title, value) {
    const item = this.shadow.querySelector(`stat-list-item[title="${title}"]`)

    if (item == null) {
      this.addItem(title, value)
    } else {
      item.value = value
    }
  }
}

class StatListItem extends HTMLElement{  
  constructor() {
    super();
  }

  set title(title) {
    this.setAttribute('title', title)
    this.innerHTML = `<span part='list-title'>${title}</span><span part='list-value'>${this.value}</span>`
  }

  set value(value) {
    this.setAttribute('value', value)
    this.innerHTML = `<span part='list-title'>${this.title}</span><span part='list-value'>${value}</span>`
  }

  get title() {
    return this.getAttribute('title')
  }

  get value() {
    return this.getAttribute('value')
  }
}

const applyTheme = async (theme) => {
  if (theme == null)
    theme = 'default'

  if (theme != 'default') {
    const cssEl = document.createElement('link')
    cssEl.setAttribute('href', `/themes/${theme}.css`)
    cssEl.setAttribute('rel', 'stylesheet')
    document.head.appendChild(cssEl)
    document.getElementById('default-theme').remove()
  }
}

const defineCustomElements = () => {
  customElements.define('stat-list', StatList);
  customElements.define('stat-list-item', StatListItem);
}

const convertToNaturalLanguage = (stat, value) => {
  let s = stat, v = value
  switch (stat) {
    case 'cfn': case 'lp': case 'mr':
      s = stat.toUpperCase()
      break
    case 'lpGain': 
      s = stat.replace('lp', 'LP ')
      if (value > 0)
        v = '+' + value
      break
    case 'mrGain': 
      s = stat.replace('mr', 'MR ')
      if (value > 0)
        v = '+' + value
      break
    case 'result':
      v = value ? 'W' : 'L'
      break
    case 'winRate':
      v = value + '%'
      break
  }

  if (stat.includes('win') && stat != 'wins') {
    s = stat.replace('win', 'Win ')
  } else if (stat.includes('total')) {
    s = stat.replace('total', 'Total ')
  } else if (stat.includes('opponent')) {
    s = stat.replace('opponent', 'Opponent ')
  }

  return { stat: s, value: v }
}

const main = () => {
  defineCustomElements()

  const searchParams = new URLSearchParams(window.location.search);
  applyTheme(searchParams.get('theme'))

  const list = document.createElement('stat-list')
  document.body.appendChild(list)

  const src = new EventSource('/stream')
  src.addEventListener('open', _ => console.log('The connection has been established'))
  src.addEventListener('error', _ => console.log('An error occurred while attempting to connect'))
  src.addEventListener('message', e => {
    const mh = JSON.parse(e.data)
    console.log('New match played: ', mh)

    for (const [stat, value] of Object.entries(mh)) {
      if (searchParams.get(stat) != 'true') 
        continue
        
      const natural = convertToNaturalLanguage(stat, value)
      list.updateItem(natural.stat, natural.value)
    }
  })
}

document.addEventListener('DOMContentLoaded', main)
