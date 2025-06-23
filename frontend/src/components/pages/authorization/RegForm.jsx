import EnvelopeIcon from '@/components/shared/icons/EnvelopeIcon';
import LockIcon from '@/components/shared/icons/LockIcon';
import UserIcon from '@/components/shared/icons/UserIcon';
import BaseHeader from '@/components/UI/BaseHeader';
import DarkInput from '@/components/UI/DarkInput';
import OrangeButton from '@/components/UI/OrangeButton';
import { routes } from '@/router';
import { Link } from 'react-router-dom';


function RegForm() {
  return (
    <>
      <BaseHeader>Давайте знакомиться</BaseHeader>

      {/* Login */}
      <div className="flex gap-3 flex-col">
        <DarkInput placeholder="Логин" type="text">
          <UserIcon className="absolute top-[25%] left-3 stroke-secondary" />
        </DarkInput>

      {/* Email */}
        <DarkInput placeholder="Email" type="email">
          <EnvelopeIcon className="absolute top-[25%] left-3 stroke-secondary" />
        </DarkInput>

      {/* Password */}
        <DarkInput placeholder="Пароль" type="password">
          <LockIcon className="absolute top-[25%] left-3 stroke-secondary" />
        </DarkInput>
      </div>
      <OrangeButton>Зарегистрироваться</OrangeButton>

      <div className="mt-14 w-[100%] text-center text-md">
        <span className="text-secondary">Есть аккаунт?</span>
        <Link
          to={{ hash: routes.public.auth.hashes.auth }}
          className="underline text-accent-1 text-xl ml-2"
        >
          Войти
        </Link>
      </div>
    </>
  );
}

export default RegForm;
