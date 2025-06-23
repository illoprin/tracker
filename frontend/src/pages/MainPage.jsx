function MainPage() {
  return (
    <>
      <div class="landing-backdrop"></div>

      <div class="flex flex-col w-screen h-screen items-center justify-center">
        <h1 class="text-center text-[16.5rem] text-accent-1 font-thin italic leading-none tracking-[-2.1rem]">
          ТРЕКЕР
        </h1>
        <h1 class="text-center text-lg-1 uppercase font-extralight">
          открыт для нового
        </h1>

        <div class="mt-10">
          <button class="btn-primary rounded-sm text-md p-2 pl-8 pr-8">
            Зарегистрироваться
          </button>
        </div>
      </div>
    </>
  );
}

export default MainPage;
