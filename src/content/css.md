
## Creating your own Browser Source theme with CSS

Create a CSS file in the themes directory and it will be available as an option in the theme selector on the output page. In the CSS there are 4 CSS selectors you can use:

- **stat-list**: the container for all of the stats
- **stat-list::part(list-item)**: the container for one stat
- **stat-list::part(list-title)**: the stat title e.g. "LP Gain", "Win Rate" ...
- **stat-list::part(list-value)**: the value of the stat

Here is *nord.css* as an example:
```css
@import url("https://fonts.googleapis.com/css2?family=Caesar+Dressing&display=swap");

body {
  background: #2e3440; /* Doesn't matter, the bg is removed in OBS anyway */
}

/* Stat list container */
stat-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  max-width: 816px;

  font-family: "Caesar Dressing", "Trebuchet MS"; /* Any font installed on the system works */
  font-size: 2rem;

  --stat-title-color: #eceff4; /* Stat title color */
  --stat-value-color: #88c0d0; /* Stat value color */
  --stat-bg-color: #4c566a; /* Stat box background color */
}

/* Stat box container */
stat-list::part(list-item) {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 0.75rem;
  white-space: nowrap;
}

/* Stat title (Wins, Losses, Win Rate, LP Gain etc.) */
stat-list::part(list-title) {}

/* Stat value */
stat-list::part(list-value) {}
```