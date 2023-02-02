import styles from "../styles/DashLayout.module.scss"
import { useRouter } from "next/router";
import Link from "next/link";

const pages = [{
  title: 'Search',
  href: '/dashboard/search',
}, {
  title: 'Saved',
  href: '/dashboard/saved',
}, {
  title: 'Timer',
  href: '/timer'
}];

export default function DashLayout({ children }: { children: JSX.Element }) {

  const router = useRouter();
  const currentRoute = router.pathname;
  return (
    <div className={styles.container}>
      <aside className={styles.sidebar}>
        {pages.map(navItem => (
          <div className={currentRoute.includes(navItem.title.toLowerCase()) ? styles.active : ''} key={navItem.title}>
            <Link href={navItem.href}>{navItem.title}</Link>
          </div>
        ))}
      </aside>
      <main>{children}</main>
    </div>
  )
}