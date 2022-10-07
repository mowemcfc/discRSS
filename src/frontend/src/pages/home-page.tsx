
import React from "react";
import { LoginButton } from "../components/login"

export const HomePage = () => {
  return (
    <section className="bg-gray-900 text-white">
      <div
        className="mx-auto max-w-screen-xl px-4 py-32 lg:flex lg:h-screen lg:items-center"
      >
        <div className="mx-auto max-w-3xl text-center">
          <h1
            className="bg-gradient-to-r from-blue-300 via-purple-500 to-blue-600 bg-clip-text text-3xl font-extrabold text-transparent sm:text-5xl"
          >
            Your RSS feed.
            <span className="sm:block"> In your own Discord Server. </span>
          </h1>

          <p className="mx-auto mt-4 max-w-xl sm:text-xl sm:leading-relaxed">
            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Nesciunt illo
            tenetur fuga ducimus numquam ea!
          </p>

          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <LoginButton />
          </div>
        </div>
      </div>
    </section>
  )
}