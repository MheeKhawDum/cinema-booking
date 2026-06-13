const db = connect("mongodb://admin:secret@mongo:27017/cinema?authSource=admin");

db.movies.deleteMany({});
db.showtimes.deleteMany({});

print("🎬 Seeding movies...");

const now = new Date();

const movies = [
  {
    _id: ObjectId("665000000000000000000001"),
    title: "Avengers: Doomsday",
    description: "The Avengers reassemble to face their greatest threat yet — a multiversal war that could erase all of existence.",
    poster_url: "https://image.tmdb.org/t/p/w500/1TRFmLaRHiRhcI7VTlGaMqjFHbE.jpg",
    duration: 150,
    created_at: now
  },
  {
    _id: ObjectId("665000000000000000000002"),
    title: "Mission: Impossible – The Final Reckoning",
    description: "Ethan Hunt races against time to stop a rogue AI from triggering global chaos in his most dangerous mission yet.",
    poster_url: "https://image.tmdb.org/t/p/w500/z53D9BbEFhJEFHKE1GW8oo0Zyph.jpg",
    duration: 169,
    created_at: now
  },
  {
    _id: ObjectId("665000000000000000000003"),
    title: "Jurassic World Rebirth",
    description: "A new expedition ventures into an uncharted island teeming with prehistoric predators, unleashing chaos beyond imagination.",
    poster_url: "https://image.tmdb.org/t/p/w500/oVfZiSbGbWljECoktbMwFpqXiDR.jpg",
    duration: 130,
    created_at: now
  },
  {
    _id: ObjectId("665000000000000000000004"),
    title: "How to Train Your Dragon",
    description: "A young Viking discovers an injured dragon and forms an unlikely friendship that challenges centuries of tradition.",
    poster_url: "https://image.tmdb.org/t/p/w500/q9EyHGSuQ0OoIqBlnHDMQRCGEHo.jpg",
    duration: 95,
    created_at: now
  },
  {
    _id: ObjectId("665000000000000000000005"),
    title: "F1",
    description: "A Formula 1 legend comes out of retirement to mentor a young rookie, facing off against rivals on the world's most dangerous circuits.",
    poster_url: "https://image.tmdb.org/t/p/w500/mGVrXeIjyYXnRFMbBmEwPxzDEV3.jpg",
    duration: 142,
    created_at: now
  }
];

db.movies.insertMany(movies);
print(`✅ Inserted ${movies.length} movies`);

print("🎟️  Seeding showtimes...");

const halls = ["Hall A", "Hall B", "Hall IMAX"];
const showtimes = [];

movies.forEach((movie, mIdx) => {
  for (let dayOffset = 0; dayOffset < 3; dayOffset++) {
    const times = ["10:00", "14:00", "18:00", "21:00"];
    const selectedTimes = times.slice(mIdx % 2, (mIdx % 2) + 3);

    selectedTimes.forEach((time, tIdx) => {
      const [hour, minute] = time.split(":").map(Number);
      const startTime = new Date(now);
      startTime.setDate(startTime.getDate() + dayOffset);
      startTime.setHours(hour, minute, 0, 0);

      showtimes.push({
        _id: new ObjectId(),
        movie_id: movie._id,
        start_time: startTime,
        hall: halls[(mIdx + tIdx) % halls.length],
        total_seats: 100,
        created_at: now
      });
    });
  }
});

db.showtimes.insertMany(showtimes);
print(`✅ Inserted ${showtimes.length} showtimes`);

db.movies.createIndex({ title: 1 });
db.showtimes.createIndex({ movie_id: 1, start_time: 1 });
db.bookings.createIndex({ showtime_id: 1, status: 1 });
db.bookings.createIndex({ user_id: 1, created_at: -1 });
