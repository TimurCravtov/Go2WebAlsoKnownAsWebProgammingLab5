## Just some website which behave predictable, can be used for the demo

### This is for content negociation

```bash
go run . -u "https://icanhazdadjoke.com/" --type application/json
```

```bash
go run . -u "https://wttr.in/Chisinau" --lang ja
```

```bash
go run . -u "https://httpbin.org/absolute-redirect/4" --max-redirects 3
```

```bash
go run . -u "https://cdn.marvel.com/content/2x/005smp_ons_cut_mob_01_6.jpg" # spiderman
go run . -u "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR1iJ3bIbJEyV-7N_yZzxbWqIAP4ANxLbYQxxJQ_xy6DA&s&ec=121585071" #silly cat
```