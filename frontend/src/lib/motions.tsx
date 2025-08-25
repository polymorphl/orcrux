export const bindVariants = {
  container: {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.2
      }
    }
  },
  item: {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.4, ease: "easeOut" as const }
    }
  },
  button: {
    hover: {
      scale: 1.02,
      transition: { duration: 0.2, ease: "easeInOut" as const }
    },
    tap: {
      scale: 0.98,
      transition: { duration: 0.1 }
    }
  },
  result: {
    hidden: { opacity: 0, height: 0 },
    visible: {
      opacity: 1,
      height: "auto",
      transition: { duration: 0.5, ease: "easeOut" as const }
    }
  }
}

export const splitFormVariants = {
  container: {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.2
      }
    }
  },
  item: {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.4, ease: "easeOut" as const }
    }
  },
  button: {
    hover: {
      scale: 1.02,
      transition: { duration: 0.2, ease: "easeInOut" as const }
    },
    tap: {
      scale: 0.98,
      transition: { duration: 0.1 }
    }
  }
}

export const splitResultVariants = {
  cardVariants: {
    hidden: { opacity: 0, y: 15 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.3, ease: "easeOut" as const }
    }
  },
  container: {
    hidden: { opacity: 0, y: 20, scale: 0.95 },
    visible: {
      opacity: 1,
      y: 0,
      scale: 1,
      transition: { duration: 0.4, ease: "easeOut" as const }
    }
  }
}